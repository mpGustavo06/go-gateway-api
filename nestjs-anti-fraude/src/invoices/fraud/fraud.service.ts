/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { Injectable } from '@nestjs/common';
import { Account, FraudReason, InvoiceStatus } from 'generated/prisma';
import { PrismaService } from 'src/prisma/prisma.service';
import { ProcessInvoiceFraudDto } from '../dto/process-invoice-fraud.dto';
import { FraudAggregateSpecification } from './specifications/fraud-aggregate.specification';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class FraudService {
  constructor(
    private prismaService: PrismaService,
    // private configService: ConfigService,
    private fraudAggregateSpec: FraudAggregateSpecification,
  ) {}

  async processInvoice(processInvoiceFraudDto: ProcessInvoiceFraudDto) {
    const { invoice_id, account_id, amount } = processInvoiceFraudDto;

    const foundInvoice = await this.prismaService.foundInvoice.findUnique({
      where: { id: invoice_id },
    });

    if (foundInvoice) {
      throw new Error('Invoice has already been processed');
    }

    const account = await this.prismaService.account.upsert({
      where: {
        id: account_id,
      },
      update: {},
      create: {
        id: account_id,
      },
    });

    const fraudResult = await this.fraudAggregateSpec.detectFraud({
      account,
      amount,
      invoiceId: invoice_id,
    });

    const invoice = await this.prismaService.foundInvoice.create({
      data: {
        id: invoice_id,
        accountId: account_id,
        amount,
        ...(fraudResult.hasFraud && {
          fraudHistory: {
            create: {
              reason: fraudResult.reason!,
              description: fraudResult.description,
            },
          },
        }),
        status: fraudResult.hasFraud
          ? InvoiceStatus.REJECTED
          : InvoiceStatus.APPROVED,
      },
    });

    return {
      invoice,
      fraudResult,
    };
  }

  // async detectFraud(data: { account: Account; amount: number }) {
  //   const { account, amount } = data;

  //   const SUSPICIOUS_VARIATION_PERCENTAGE =
  //     this.configService.getOrThrow<number>('SUSPICIOUS_VARIATION_PERCENTAGE');
  //   const INVOICES_HISTORY_COUNT = this.configService.getOrThrow<number>(
  //     'INVOICES_HISTORY_COUNT',
  //   );
  //   const SUSPICIOUS_INVOICES_COUNT = this.configService.getOrThrow<number>(
  //     'SUSPICIOUS_INVOICES_COUNT',
  //   );
  //   const SUSPICIOUS_TIMEFRAME_HOURS = this.configService.getOrThrow<number>(
  //     'SUSPICIOUS_TIMEFRAME_HOURS',
  //   );

  //   if (account.isSuspicious) {
  //     return {
  //       hasFraud: true,
  //       reason: FraudReason.SUSPICIOUS_ACCOUNT,
  //       description: 'Suspicious account',
  //     };
  //   }

  //   const previousInvoices = await this.prismaService.foundInvoice.findMany({
  //     where: {
  //       accountId: account.id,
  //     },
  //     orderBy: {
  //       createdAt: 'desc',
  //     },
  //     take: INVOICES_HISTORY_COUNT,
  //   });

  //   if (previousInvoices.length) {
  //     const totalAmount = previousInvoices.reduce((acc, foundInvoice) => {
  //       return acc + foundInvoice.amount;
  //     }, 0);

  //     const averageAmount = totalAmount / previousInvoices.length;

  //     if (
  //       amount >
  //       averageAmount * (1 + SUSPICIOUS_VARIATION_PERCENTAGE / 100) +
  //         averageAmount
  //     ) {
  //       return {
  //         hasFraud: true,
  //         reason: FraudReason.UNUSUAL_PATTERN,
  //         description: `Amount ${amount} is too high compared to average amount ${averageAmount}`,
  //       };
  //     }
  //   }

  //   const recentDate = new Date();
  //   recentDate.setHours(recentDate.getHours() - SUSPICIOUS_TIMEFRAME_HOURS);

  //   const recentInvoices = await this.prismaService.foundInvoice.findMany({
  //     where: {
  //       accountId: account.id,
  //       createdAt: {
  //         gte: recentDate,
  //       },
  //     },
  //   });

  //   if (recentInvoices.length >= SUSPICIOUS_INVOICES_COUNT) {
  //     return {
  //       hasFraud: true,
  //       reason: FraudReason.FREQUENT_HIGH_VALUE,
  //       description: `Account ${account.id} has made ${SUSPICIOUS_INVOICES_COUNT} invoices in the last ${SUSPICIOUS_TIMEFRAME_HOURS} hours`,
  //     };
  //   }

  //   return {
  //     hasFraud: false,
  //   };
  // }
}
